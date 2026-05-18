import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { format, parseISO } from 'date-fns'
import {
  createHall as createHallApi,
  createMovie,
  createSession as createSessionApi,
  deleteMovie,
  updateMovie,
} from '../../api/admin'
import { getMovies } from '../../api/movies'
import type { Movie } from '../../types'
import { MovieForm, type MovieFormValues } from '../../components/admin/MovieForm'
import { HallForm } from '../../components/admin/HallForm'
import { SessionForm, type SessionFormValues } from '../../components/admin/SessionForm'
import { Button } from '../../components/ui/Button'
import { ErrorBanner } from '../../components/ui/ErrorBanner'
import { Spinner } from '../../components/ui/Spinner'
import { getErrorMessage } from '../../utils/errorHandler'

export function AdminMoviesPage() {
  const queryClient = useQueryClient()
  const [showAdd, setShowAdd] = useState(false)
  const [editing, setEditing] = useState<Movie | null>(null)
  const [halls, setHalls] = useState<{ id: string; name: string; city?: string }[]>([])

  const moviesQuery = useQuery({
    queryKey: ['movies', 'admin'],
    queryFn: () => getMovies({ limit: 200 }),
  })

  const createMovieMut = useMutation({
    mutationFn: (body: MovieFormValues) =>
      createMovie({
        title: body.title,
        description: body.description ?? '',
        genre: body.genre,
        duration: body.duration,
        rating: body.rating,
        age_rating: body.age_rating,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['movies'] })
      setShowAdd(false)
    },
  })

  const updateMovieMut = useMutation({
    mutationFn: (vars: { id: string; body: Partial<MovieFormValues> }) =>
      updateMovie(vars.id, {
        title: vars.body.title,
        description: vars.body.description,
        genre: vars.body.genre,
        duration: vars.body.duration,
        rating: vars.body.rating,
        age_rating: vars.body.age_rating,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['movies'] })
      setEditing(null)
    },
  })

  const removeMovie = useMutation({
    mutationFn: (id: string) => deleteMovie(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['movies'] })
    },
  })

  const createHallMut = useMutation({
    mutationFn: (values: {
      name: string
      capacity: number
      city: string
      cinema_name: string
    }) =>
      createHallApi(values.name, values.capacity, values.city, values.cinema_name),
    onSuccess: (hall) => {
      setHalls((prev) => [
        ...prev,
        {
          id: hall.id,
          name: `${hall.cinema_name || 'Cinema'} · ${hall.name}`,
          city: hall.city,
        },
      ])
    },
  })

  const createSessionMut = useMutation({
    mutationFn: (body: SessionFormValues) =>
      createSessionApi(
        body.movie_id,
        body.hall_id,
        body.start_time,
        body.price,
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  const movies = moviesQuery.data ?? []

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-section font-light text-white">Movies</h1>
        <Button type="button" variant="secondary" onClick={() => setShowAdd((v) => !v)}>
          {showAdd ? 'Close form' : 'Add movie'}
        </Button>
      </div>
      {showAdd ? (
        <MovieForm
          submitLabel="Create movie"
          isSubmitting={createMovieMut.isPending}
          onSubmit={async (values) => {
            await createMovieMut.mutateAsync(values)
          }}
          onCancel={() => setShowAdd(false)}
        />
      ) : null}
      {createMovieMut.isError ? (
        <ErrorBanner message={getErrorMessage(createMovieMut.error)} />
      ) : null}
      {moviesQuery.isLoading ? (
        <div className="flex justify-center py-12">
          <Spinner className="h-6 w-6" />
        </div>
      ) : moviesQuery.isError ? (
        <ErrorBanner message={getErrorMessage(moviesQuery.error)} />
      ) : (
        <div className="overflow-x-auto rounded-lg border border-border bg-card">
          <table className="min-w-full divide-y divide-border text-left text-body">
            <thead className="text-[11px] font-medium text-muted">
              <tr>
                <th className="px-4 py-2">Title</th>
                <th className="px-4 py-2">Genre</th>
                <th className="px-4 py-2">Duration</th>
                <th className="px-4 py-2">Rating</th>
                <th className="px-4 py-2">Created</th>
                <th className="px-4 py-2">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-white">
              {movies.map((m) => (
                <tr key={m.id}>
                  <td className="px-4 py-2 font-medium">{m.title}</td>
                  <td className="px-4 py-2">{m.genre}</td>
                  <td className="px-4 py-2">{m.duration} min</td>
                  <td className="px-4 py-2">{m.rating.toFixed(1)}</td>
                  <td className="px-4 py-2 text-muted">
                    {m.created_at
                      ? format(parseISO(m.created_at), 'PP')
                      : '—'}
                  </td>
                  <td className="space-x-2 px-4 py-2">
                    <Button
                      type="button"
                      variant="secondary"
                      className="text-xs"
                      onClick={() => setEditing(m)}
                    >
                      Edit
                    </Button>
                    <Button
                      type="button"
                      variant="secondary"
                      className="text-xs"
                      disabled={removeMovie.isPending}
                      onClick={() => {
                        if (window.confirm('Delete this movie?')) {
                          void removeMovie.mutateAsync(m.id)
                        }
                      }}
                    >
                      Delete
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {editing ? (
        <div className="rounded-lg border border-border bg-card2 p-4">
          <p className="mb-3 text-body font-medium text-white">Edit movie</p>
          <MovieForm
            initial={editing}
            submitLabel="Save changes"
            isSubmitting={updateMovieMut.isPending}
            onSubmit={async (values) => {
              await updateMovieMut.mutateAsync({ id: editing.id, body: values })
            }}
            onCancel={() => setEditing(null)}
          />
          {updateMovieMut.isError ? (
            <div className="mt-3">
              <ErrorBanner message={getErrorMessage(updateMovieMut.error)} />
            </div>
          ) : null}
        </div>
      ) : null}
      <div className="grid gap-4 lg:grid-cols-2">
        <div className="flex flex-col gap-2">
          <HallForm
            isSubmitting={createHallMut.isPending}
            onSubmit={async (values) => {
              await createHallMut.mutateAsync(values)
            }}
          />
          {createHallMut.isError ? (
            <ErrorBanner message={getErrorMessage(createHallMut.error)} />
          ) : null}
        </div>
        <div className="flex flex-col gap-2">
          {movies.length > 0 ? (
            <SessionForm
              movies={movies}
              hallOptions={halls}
              isSubmitting={createSessionMut.isPending}
              onSubmit={async (values) => {
                await createSessionMut.mutateAsync({
                  ...values,
                  start_time: new Date(values.start_time).toISOString(),
                })
              }}
            />
          ) : (
            <p className="rounded-lg border border-border bg-card p-4 text-body text-muted">
              Add a movie before creating sessions.
            </p>
          )}
          {createSessionMut.isError ? (
            <ErrorBanner message={getErrorMessage(createSessionMut.error)} />
          ) : null}
        </div>
      </div>
      {halls.length === 0 ? (
        <p className="text-body text-muted">
          Create a hall before scheduling sessions.
        </p>
      ) : null}
    </div>
  )
}
