package taskimpl

// NOTE: must import all packages that contain tasks so their init functions are called.
import (
	_ "nebula/pkg/taskimpl/foots"
	_ "nebula/pkg/taskimpl/shop"
	_ "nebula/pkg/taskimpl/ys"
)
