package main

import (
	"os"

	"github.com/dewzzjr/amartha/reconcile"
)

func main() {
	reconcile.Reconcile(os.Args)
}
