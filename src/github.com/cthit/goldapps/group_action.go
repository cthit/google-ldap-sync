package goldapps

import (
	"fmt"
	"bytes"
)

// Set of action, to be performed on a set of groups
type GroupActions struct {
	Updates   []GroupUpdate
	Additions []Group
	Deletions []Group
}
func (actions GroupActions) Amount() int {
	return len(actions.Additions) + len(actions.Deletions) + len(actions.Updates)
}

// Set of actions that could not be performed with accompanying errors
type GroupActionErrors struct {
	Updates   []GroupUpdateError
	Additions []GroupAddOrDelError
	Deletions []GroupAddOrDelError
}
type GroupUpdateError struct {
	Action GroupUpdate
	Error  error
}
type GroupAddOrDelError struct {
	Action Group
	Error  error
}
func (actions GroupActionErrors) Amount() int {
	return len(actions.Additions) + len(actions.Deletions) + len(actions.Updates)
}
func (actions GroupActionErrors) String() string {
	builder := bytes.Buffer{}
	for _,deletion := range actions.Deletions  {
		builder.WriteString(fmt.Sprintf("Deletion of group \"%s\" failed with error %s\n", deletion.Action.Email, deletion.Error.Error()))
	}
	for _,update := range actions.Updates  {
		builder.WriteString(fmt.Sprintf("Update of group \"%s\" failed with error %s\n", update.Action.After.Email, update.Error.Error()))
	}
	for _,addition := range actions.Additions  {
		builder.WriteString(fmt.Sprintf("Addition of group \"%s\" failed with error %s\n", addition.Action.Email, addition.Error.Error()))
	}
	return builder.String()
}

// Data struct representing how a group looks not and how it should look after an update
// Allows for efficient updates as application doesn't have to re-upload whole group
type GroupUpdate struct {
	Before Group
	After  Group
}

// Commits a set of actions to a service.
// Returns all actions performed and a error if not all actions could be performed for some reason.
func (actions GroupActions) Commit(service UpdateService) GroupActionErrors {

	errors := GroupActionErrors{}

	if len(actions.Deletions) > 0 {
		fmt.Println("(Groups) Performing deletions")
		printProgress(0, len(actions.Deletions), 0)
		for deletionsIndex, group := range actions.Deletions {
			err := service.DeleteGroup(group)
			if err != nil {
				// Save error
				errors.Deletions = append(errors.Deletions, GroupAddOrDelError{Action: group, Error: err})
			}
			printProgress(deletionsIndex+1, len(actions.Deletions), len(errors.Deletions))
		}
	}

	if len(actions.Updates) > 0 {
		fmt.Println("(Groups) Performing updates")
		printProgress(0, len(actions.Updates), 0)
		for updatesIndex, update := range actions.Updates {
			err := service.UpdateGroup(update)
			if err != nil {
				// Save error
				errors.Updates = append(errors.Updates, GroupUpdateError{Action: update, Error: err})
			}
			printProgress(updatesIndex+1, len(actions.Updates), len(errors.Updates))
		}
	}

	if len(actions.Additions) > 0 {
		fmt.Println("(Groups) Performing additions")
		printProgress(0, len(actions.Additions), 0)
		for additionsIndex, group := range actions.Additions {
			err := service.AddGroup(group)
			if err != nil {
				// Save error
				errors.Additions = append(errors.Additions, GroupAddOrDelError{Action: group, Error: err})
			}
			printProgress(additionsIndex+1, len(actions.Additions), len(errors.Additions))
		}
	}

	return errors
}

// Determines actions required to make the "old" group list look as the "new" group list.
// Returns a list with those actions.
func GroupActionsRequired(old []Group, new []Group) GroupActions {
	requiredActions := GroupActions{}

	for _, newGroup := range new {
		exists := false
		for _, oldGroup := range old {
			// identify by Email
			if newGroup.Email == oldGroup.Email {
				// Groups exists
				exists = true
				// check if group has to be updates
				if !newGroup.equals(oldGroup) {
					// Add group update
					requiredActions.Updates = append(requiredActions.Updates, GroupUpdate{
						Before: oldGroup,
						After:  newGroup,
					})
				}
				break
			}
		}

		// Add group creation action if group doesn't exist
		if !exists {
			requiredActions.Additions = append(requiredActions.Additions, newGroup)
		}
	}

	for _, oldGroup := range old {
		// check if group should be removed
		removed := true
		for _, newGroup := range new {
			if oldGroup.Email == newGroup.Email {
				removed = false
				break
			}
		}

		if removed {
			// Add group deletion action
			requiredActions.Deletions = append(requiredActions.Deletions, oldGroup)
		}
	}

	return requiredActions
}
