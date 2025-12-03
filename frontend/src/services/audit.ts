import api from './api';
import type { GetAuditLogsParams, AuditLogsResponse } from '../types/audit';

export async function getAuditLogs(params: GetAuditLogsParams = {}): Promise<AuditLogsResponse> {
  const resp = await api.get('/admin/audit-logs', { params });
  return resp.data as AuditLogsResponse;
}

export default { getAuditLogs };
