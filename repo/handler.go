package repo

import "context"

func WithinTransaction(ctx context.Context, s App, tFunc func(ctx context.Context) error) error {
	var isNewTx bool
	tx := ExtractTx(ctx)
	// there is no tx inject to ctx, so new one
	if tx == nil {
		tx = s.Begin()
		defer tx.Rollback()
		isNewTx = true
	}

	// run callback
	if err := tFunc(InjectTx(ctx, tx)); err != nil {
		return err
	}

	// tx extract from ctx don't need commit
	if isNewTx {
		return tx.Commit()
	}
	return nil
}
