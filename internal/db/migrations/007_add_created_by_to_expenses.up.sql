ALTER TABLE expenses
  ADD COLUMN created_by UUID REFERENCES users(id);
