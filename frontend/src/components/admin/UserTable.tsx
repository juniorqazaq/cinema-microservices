import { format, parseISO } from 'date-fns'
import type { User } from '../../types'
import { Badge } from '../ui/Badge'
import { Button } from '../ui/Button'
import { formatPrice } from '../../utils/format'

interface UserTableProps {
  users: User[]
  onBan: (userId: string) => void
  onRoleChange: (userId: string, role: 'user' | 'admin') => void
  banningId: string | null
  roleUpdatingId: string | null
}

export function UserTable({
  users,
  onBan,
  onRoleChange,
  banningId,
  roleUpdatingId,
}: UserTableProps) {
  return (
    <div className="overflow-x-auto rounded-lg border border-border bg-card">
      <table className="min-w-full divide-y divide-border text-left text-body">
        <thead className="text-[11px] font-medium text-muted">
          <tr>
            <th className="px-4 py-2">Email</th>
            <th className="px-4 py-2">Name</th>
            <th className="px-4 py-2">Role</th>
            <th className="px-4 py-2">Balance</th>
            <th className="px-4 py-2">Status</th>
            <th className="px-4 py-2">Registered</th>
            <th className="px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border text-white">
          {users.map((u) => (
            <tr key={u.id}>
              <td className="px-4 py-2">{u.email}</td>
              <td className="px-4 py-2 text-muted">{u.full_name || '—'}</td>
              <td className="px-4 py-2">
                <Badge tone={u.role === 'admin' ? 'accent' : 'neutral'}>
                  {u.role}
                </Badge>
              </td>
              <td className="px-4 py-2">{formatPrice(u.balance ?? 0)}</td>
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
              <td className="space-x-2 px-4 py-2">
                {!u.is_banned ? (
                  <Button
                    type="button"
                    variant="secondary"
                    className="text-xs"
                    disabled={banningId === u.id}
                    onClick={() => onBan(u.id)}
                  >
                    Ban
                  </Button>
                ) : null}
                <Button
                  type="button"
                  variant="secondary"
                  className="text-xs"
                  disabled={roleUpdatingId === u.id}
                  onClick={() =>
                    onRoleChange(u.id, u.role === 'admin' ? 'user' : 'admin')
                  }
                >
                  {u.role === 'admin' ? 'Make user' : 'Make admin'}
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
