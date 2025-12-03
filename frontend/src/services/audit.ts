import api from './api';
import type { GetAuditLogsParams, AuditLogsResponse, AuditLog, AuditLogSummary } from '../types/audit';

export async function getAuditLogs(params: GetAuditLogsParams = {}): Promise<AuditLogsResponse> {
  const resp = await api.get('/admin/audit-logs', { params });
  return resp.data as AuditLogsResponse;
}

export async function getAuditLog(id: string | number): Promise<AuditLog> {
  const resp = await api.get(`/admin/audit-logs/${id}`);
  return resp.data as AuditLog;
}

export async function getAuditLogDetailed(id: string | number): Promise<AuditLog> {
  const resp = await api.get(`/admin/audit-logs/${id}/detailed`);
  return resp.data as AuditLog;
}

export async function getAuditLogsSummary(): Promise<AuditLogSummary> {
  const resp = await api.get('/admin/audit-logs/summary');
  return resp.data as AuditLogSummary;
}

export default { getAuditLogs, getAuditLog, getAuditLogDetailed, getAuditLogsSummary };
