package repo

import "context"

type txKey struct{}

// injectTx injects transaction to context
func InjectTx(ctx context.Context, tx App) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func ExtractTx(ctx context.Context) App {
	if tx, ok := ctx.Value(txKey{}).(App); ok {
		return tx
	}
	return nil
}

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
