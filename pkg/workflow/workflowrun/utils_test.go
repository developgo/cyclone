package workflowrun

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
)

func TestGetResourceVolumeName(t *testing.T) {
	assert.Equal(t, "rsc-git", GetResourceVolumeName("git"))
}

func TestResolveStatus(t *testing.T) {
	latest := &v1alpha1.Status{
		Status: v1alpha1.StatusCompleted,
	}
	update := &v1alpha1.Status{
		Status: v1alpha1.StatusRunning,
	}
	expected := &v1alpha1.Status{
		Status: v1alpha1.StatusCompleted,
	}
	result := resolveStatus(latest, update)
	assert.Equal(t, expected, result)

	latest = &v1alpha1.Status{
		Status: v1alpha1.StatusRunning,
	}
	update = &v1alpha1.Status{
		Status: v1alpha1.StatusCompleted,
	}
	expected = &v1alpha1.Status{
		Status: v1alpha1.StatusCompleted,
	}
	result = resolveStatus(latest, update)
	assert.Equal(t, expected, result)

	latest = &v1alpha1.Status{
		Status: v1alpha1.StatusError,
	}
	update = &v1alpha1.Status{
		Status: v1alpha1.StatusCompleted,
	}
	expected = &v1alpha1.Status{
		Status: v1alpha1.StatusError,
	}
	result = resolveStatus(latest, update)
	assert.Equal(t, expected, result)

	now := metav1.Time{Time: time.Now()}
	old := metav1.Time{Time: time.Now().Add(-time.Second * 10)}
	latest = &v1alpha1.Status{
		Status:             v1alpha1.StatusRunning,
		LastTransitionTime: now,
	}
	update = &v1alpha1.Status{
		Status:             v1alpha1.StatusRunning,
		LastTransitionTime: old,
	}
	expected = &v1alpha1.Status{
		Status:             v1alpha1.StatusRunning,
		LastTransitionTime: now,
	}
	result = resolveStatus(latest, update)
	assert.Equal(t, expected, result)
}

func TestNextStages(t *testing.T) {
	wf := &v1alpha1.Workflow{
		Spec: v1alpha1.WorkflowSpec{
			Stages: []v1alpha1.StageItem{
				{
					Name: "A",
				},
				{
					Name:    "B",
					Depends: []string{"A"},
				},
				{
					Name: "C",
				},
			},
		},
	}
	wfr := &v1alpha1.WorkflowRun{
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Status: v1alpha1.StatusCompleted},
				},
			},
		},
	}
	expected := []string{"B", "C"}
	nexts := NextStages(wf, wfr)
	assert.Equal(t, expected, nexts)

	wf = &v1alpha1.Workflow{
		Spec: v1alpha1.WorkflowSpec{
			Stages: []v1alpha1.StageItem{
				{
					Name: "A",
				},
				{
					Name:    "B",
					Depends: []string{"A"},
				},
				{
					Name: "C",
				},
			},
		},
	}
	wfr = &v1alpha1.WorkflowRun{
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Status: v1alpha1.StatusError},
				},
			},
		},
	}
	expected = []string{"C"}
	nexts = NextStages(wf, wfr)
	assert.Equal(t, expected, nexts)

	wf = &v1alpha1.Workflow{
		Spec: v1alpha1.WorkflowSpec{
			Stages: []v1alpha1.StageItem{
				{
					Name: "A",
				},
				{
					Name:    "B",
					Depends: []string{"A"},
				},
				{
					Name: "C",
				},
			},
		},
	}
	wfr = &v1alpha1.WorkflowRun{
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Status: v1alpha1.StatusRunning},
				},
			},
		},
	}
	expected = []string{"C"}
	nexts = NextStages(wf, wfr)
	assert.Equal(t, expected, nexts)
}

func TestStaticStatus(t *testing.T) {
	now := metav1.Time{Time: time.Now()}
	zero := metav1.Time{Time: time.Unix(0, 0)}
	status := &v1alpha1.WorkflowRunStatus{
		Stages: map[string]*v1alpha1.StageStatus{
			"A": {
				Status: v1alpha1.Status{
					Status:             v1alpha1.StatusRunning,
					LastTransitionTime: now,
				},
			},
		},
		Overall: v1alpha1.Status{
			Status:             v1alpha1.StatusRunning,
			LastTransitionTime: now,
		},
	}
	actual := staticStatus(status)
	expected := &v1alpha1.WorkflowRunStatus{
		Stages: map[string]*v1alpha1.StageStatus{
			"A": {
				Status: v1alpha1.Status{
					Status:             v1alpha1.StatusRunning,
					LastTransitionTime: zero,
				},
			},
		},
		Overall: v1alpha1.Status{
			Status:             v1alpha1.StatusRunning,
			LastTransitionTime: zero,
		},
	}
	assert.Equal(t, expected, actual)
}

func TestString(t *testing.T) {
	item := workflowRunItem{
		name:      "test",
		namespace: "default",
	}
	assert.Equal(t, "default:test", item.String())
}
