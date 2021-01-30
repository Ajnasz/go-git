package object

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type commitDifferenceIterator struct {
	except     map[plumbing.Hash]struct{}
	sourceIter CommitIter
	start      *Commit
}

// NewCommitDifferenceIterFromIter returns a commit iter that walkd the commit
// history like WalkCommitHistory but filters out the commits which are not in
// the seen hash
func NewCommitDifferenceIterFromIter(except map[plumbing.Hash]struct{}, commitIter CommitIter) CommitIter {
	iterator := new(commitDifferenceIterator)
	iterator.except = except
	iterator.sourceIter = commitIter

	return iterator
}

func (c *commitDifferenceIterator) Next() (*Commit, error) {
	for {
		commit, err := c.sourceIter.Next()

		if err != nil {
			return nil, err
		}

		if _, ok := c.except[commit.Hash]; ok {
			continue
		}

		return commit, nil
	}
}

func (c *commitDifferenceIterator) ForEach(cb func(*Commit) error) error {
	for {
		commit, nextErr := c.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return nextErr
		}
		err := cb(commit)
		if err == storer.ErrStop {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (c *commitDifferenceIterator) Close() {
	c.sourceIter.Close()
}
