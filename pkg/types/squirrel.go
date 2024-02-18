package types

import sq "github.com/Masterminds/squirrel"

type SquirrelBuilder sq.SelectBuilder || sq.UpdateBuilder || sq.DeleteBuilder
