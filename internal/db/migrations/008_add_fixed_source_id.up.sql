ALTER TABLE expenses ADD COLUMN fixed_source_id UUID REFERENCES expenses(id) ON DELETE SET NULL;
CREATE INDEX idx_expenses_fixed_source_id ON expenses(fixed_source_id);
