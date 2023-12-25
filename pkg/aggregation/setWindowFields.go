package aggregation

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type SetWindowFieldsInput struct {
	Documents   *[]bson.M
	PartitionBy string
	SortBy      string
	Asc         bool
	Output      string
}

type SetWindowFieldsOutput struct {
	Result *[]bson.M
}

type SetWindowFieldsCommand struct {
	In  *SetWindowFieldsInput
	Out *SetWindowFieldsOutput
}

func NewSetWindowFieldsCommand(in *SetWindowFieldsInput) *SetWindowFieldsCommand {
	return &SetWindowFieldsCommand{In: in}
}

// This function emulates the setWindowFields function in MongoDB
// Execute executes the SetWindowFieldsCommand.
// It groups the documents by the specified partitionBy field,
// sorts them by the specified sortBy field, and sets the window fields.
// The window fields are set based on the order of the documents within each group.
// Finally, it assigns the result to the Out field of the command.
func (c *SetWindowFieldsCommand) Execute(ctx context.Context) error {
	var result []bson.M
	// If documents is empty, return empty result
	if len(*c.In.Documents) == 0 {
		c.Out = &SetWindowFieldsOutput{
			Result: &result,
		}
		return nil
	}
	// Group by partitionBy
	groupByPartitionBy := make(map[string][]bson.M)
	for _, document := range *c.In.Documents {
		partitionByValue := document[c.In.PartitionBy]
		groupByPartitionBy[partitionByValue.(string)] = append(groupByPartitionBy[partitionByValue.(string)], document)
	}
	// Sort by sortBy
	for _, documents := range groupByPartitionBy {
		quickSort(documents, c.In.SortBy, c.In.Asc)
	}
	// Set window fields
	for _, documents := range groupByPartitionBy {
		for i, document := range documents {
			document[c.In.Output] = i + 1
			result = append(result, document)
		}
	}
	c.Out = &SetWindowFieldsOutput{
		Result: &result,
	}
	return nil
}

func quickSort(documents []bson.M, sortBy string, asc bool) {
	if len(documents) <= 1 {
		return
	}

	pivot := documents[len(documents)/2][sortBy].(int)
	left := 0
	right := len(documents) - 1

	for left <= right {
		if asc {
			for documents[left][sortBy].(int) < pivot {
				left++
			}
			for documents[right][sortBy].(int) > pivot {
				right--
			}
		} else {
			for documents[left][sortBy].(int) > pivot {
				left++
			}
			for documents[right][sortBy].(int) < pivot {
				right--
			}
		}

		if left <= right {
			documents[left], documents[right] = documents[right], documents[left]
			left++
			right--
		}
	}

	quickSort(documents[:left], sortBy, asc)
	quickSort(documents[left:], sortBy, asc)
}
