import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { banUser } from '../../api/admin'
import { useAdminUsers } from '../../hooks/useBooking'
import { UserTable } from '../../components/admin/UserTable'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import { getErrorMessage } from '../../utils/errorHandler'

export function AdminUsersPage() {
  const queryClient = useQueryClient()
  const [banningId, setBanningId] = useState<string | null>(null)
  const usersQuery = useAdminUsers(1, 50)

  const banMutation = useMutation({
    mutationFn: (userId: string) => banUser(userId),
    onMutate: async (userId) => {
      setBanningId(userId)
    },
    onSettled: () => {
      setBanningId(null)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] })
    },
  })

  if (usersQuery.isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner className="h-6 w-6" />
      </div>
    )
  }

  if (usersQuery.isError) {
    return <ErrorBanner message={getErrorMessage(usersQuery.error)} />
  }

  const users = usersQuery.data ?? []

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-section font-light text-white">Users</h1>
      {banMutation.isError ? (
        <ErrorBanner message={getErrorMessage(banMutation.error)} />
      ) : null}
      <UserTable
        users={users}
        banningId={banningId}
        onBan={(id) => {
          void banMutation.mutateAsync(id)
        }}
      />
    </div>
  )
}
