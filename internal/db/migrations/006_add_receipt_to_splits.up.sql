ALTER TABLE expense_splits
  ADD COLUMN receipt_data TEXT,
  ADD COLUMN receipt_type VARCHAR(50),
  ADD COLUMN receipt_file_name VARCHAR(255);
