CREATE INDEX IF NOT EXISTS idx_expenses_group_competence
  ON expenses(group_id, competence_date DESC, created_at DESC);
