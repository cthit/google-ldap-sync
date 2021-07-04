package gamma

import (
	"testing"

	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/magiconair/properties/assert"
)

func TestShouldHaveMail(t *testing.T) {
	userA, userB := FKITUser{
		Gdpr: true,
	}, FKITUser{
		Gdpr: false,
	}

	groupA, groupB, groupC, groupD := FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       false,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "BOARD",
		},
		Active:       false,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "BOARD",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}

	assert.Equal(t, shouldHaveMail(&groupA, &userA), false)
	assert.Equal(t, shouldHaveMail(&groupA, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupB, &userA), true)
	assert.Equal(t, shouldHaveMail(&groupB, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupC, &userA), false)
	assert.Equal(t, shouldHaveMail(&groupC, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupD, &userA), true)
	assert.Equal(t, shouldHaveMail(&groupD, &userB), false)
}

func TestGetMembers(t *testing.T) {
	userA, userB := FKITUser{
		Cid:   "usera",
		Email: "usera@gmail.com",
		Gdpr:  true,
	}, FKITUser{
		Cid:   "userb",
		Email: "userb@gmail.com",
		Gdpr:  false,
	}

	groupA, groupB := FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "SOCIETY",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB, userB},
	}

	assert.Equal(t, getMembers(&groupA), []string{"usera@chalmers.it"})
	assert.Equal(t, getMembers(&groupB), []string{"usera@gmail.com", "userb@gmail.com"})
}

func TestGetGroups(t *testing.T) {
	groups := []FKITGroup{
		{
			Active: true,
			SuperGroup: FKITSuperGroup{
				Type:  "COMMITTEE",
				Email: "digit@chalmers.it",
			},
		},
	}

	want := []model.Group{{"digit@chalmers.it", "COMMITTEE", []string{"ita.styrit@chalmers.it"}, nil, false},
		{"grupper@chalmers.it", "", []string{"digit@chalmers.it"}, nil, false},
		{"kommitteer@chalmers.it", "", []string{"digit@chalmers.it"}, nil, false}}

	assert.Equal(t, getGroups(groups), want)
}

func TestGetPostMails(t *testing.T) {
	userA := FKITUser{
		Cid: "usera",
		Post: Post{
			EmailPrefix: "ordf",
		},
		Gdpr:  true,
		Email: "usera@gamil.com",
	}

	groupA := FKITGroup{
		Active: true,
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
			Email: "supergroup@chalmers.it",
		},
		GroupMembers: []FKITUser{userA},
	}

	want := []model.Group{{"ordf.supergroup@chalmers.it", "", []string{"usera@chalmers.it"}, nil, false},
		{"ordforanden@chalmers.it", "", []string{"ordf.supergroup@chalmers.it"}, nil, false},
		{"ordforanden.kommitteer@chalmers.it", "", []string{"ordf.supergroup@chalmers.it"}, nil, false},
		{"kassorer@chalmers.it", "", []string{}, nil, false},
		{"kassorer.kommitteer@chalmers.it", "", []string{}, nil, false}}

	assert.Equal(t, getPostMails([]FKITGroup{groupA}), want)
}
