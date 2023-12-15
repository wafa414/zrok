package sync

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func Synchronize(src, dst Target) error {
	srcTree, err := src.Inventory()
	if err != nil {
		return errors.Wrap(err, "error creating source inventory")
	}

	dstTree, err := dst.Inventory()
	if err != nil {
		return errors.Wrap(err, "error creating destination inventory")
	}

	dstIndex := make(map[string]*Object)
	for _, f := range dstTree {
		dstIndex[f.Path] = f
	}

	var copyList []*Object
	for _, srcF := range srcTree {
		if dstF, found := dstIndex[srcF.Path]; found {
			if dstF.Size != srcF.Size || dstF.Modified.UTC() != srcF.Modified.UTC() {
				copyList = append(copyList, srcF)
			}
		} else {
			copyList = append(copyList, srcF)
		}
	}

	for _, copyPath := range copyList {
		ss, err := src.ReadStream(copyPath.Path)
		if err != nil {
			return err
		}
		if err := dst.WriteStream(copyPath.Path, ss, os.ModePerm); err != nil {
			return err
		}
		if err := dst.SetModificationTime(copyPath.Path, copyPath.Modified); err != nil {
			return err
		}
		logrus.Infof("=> %v", copyPath.Path)
	}

	return nil
}
