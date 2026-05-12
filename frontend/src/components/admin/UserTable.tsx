import { format, parseISO } from 'date-fns'
import type { User } from '../../types'
import { Badge } from '../ui/Badge'
import { Button } from '../ui/Button'

interface UserTableProps {
  users: User[]
  onBan: (userId: string) => void
  banningId: string | null
}

export function UserTable({ users, onBan, banningId }: UserTableProps) {
  return (
    <div className="overflow-x-auto rounded-lg border border-border bg-card">
      <table className="min-w-full divide-y divide-border text-left text-body">
        <thead className="text-[11px] font-medium text-muted">
          <tr>
            <th className="px-4 py-2">Email</th>
            <th className="px-4 py-2">Role</th>
            <th className="px-4 py-2">Status</th>
            <th className="px-4 py-2">Created</th>
            <th className="px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border text-white">
          {users.map((u) => (
            <tr key={u.id}>
              <td className="px-4 py-2">{u.email}</td>
              <td className="px-4 py-2">
                <Badge tone={u.role === 'admin' ? 'accent' : 'neutral'}>
                  {u.role}
                </Badge>
              </td>
              <td className="px-4 py-2">
                {u.is_banned ? (
                  <Badge tone="danger">Banned</Badge>
                ) : (
                  <span className="text-muted">Active</span>
                )}
              </td>
              <td className="px-4 py-2 text-muted">
                {u.created_at ? format(parseISO(u.created_at), 'PP') : '—'}
              </td>
              <td className="px-4 py-2">
                {!u.is_banned ? (
                  <Button
                    type="button"
                    variant="secondary"
                    className="text-sm"
                    disabled={banningId === u.id}
                    onClick={() => onBan(u.id)}
                  >
                    Ban
                  </Button>
                ) : null}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
