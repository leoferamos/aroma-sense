export type AuditActor = {
  id: number;
  public_id?: string;
  display_name?: string;
  email?: string;
  role?: string;
};

export type AuditLog = {
  id: number;
  public_id?: string;
  action: string;
  actor_id?: number;
  actor?: AuditActor | null;
  user_id?: number;
  user?: AuditActor | null;
  resource?: string;
  resource_id?: string;
  severity?: string;
  timestamp?: string;
  created_at?: string;
  details?: Record<string, unknown> | null;
  old_values?: Record<string, unknown> | null;
  new_values?: Record<string, unknown> | null;
  compliance?: string;
};

export type GetAuditLogsParams = {
  user_id?: number;
  actor_id?: number;
  action?: string;
  resource?: string;
  resource_id?: string;
  start_date?: string;
  end_date?: string;
  limit?: number;
  offset?: number;
};

export type AuditLogsResponse = {
  audit_logs: AuditLog[];
  limit: number;
  offset: number;
  total: number;
};
