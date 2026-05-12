import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { Movie } from '../../types'
import { MOVIE_GENRES } from '../../constants/genres'
import { Button } from '../ui/Button'
import { Input, Select } from '../ui/Input'
import { Spinner } from '../ui/Spinner'

const schema = z.object({
  title: z.string().min(1, 'Title is required'),
  description: z.string().optional(),
  genre: z.string().min(1, 'Genre is required'),
  duration: z.number().min(1, 'Duration must be at least 1'),
  rating: z.number().min(0).max(10),
})

export type MovieFormValues = z.infer<typeof schema>

interface MovieFormProps {
  initial?: Movie | null
  onSubmit: (values: MovieFormValues) => Promise<void>
  onCancel?: () => void
  submitLabel?: string
  isSubmitting?: boolean
}

export function MovieForm({
  initial,
  onSubmit,
  onCancel,
  submitLabel = 'Save',
  isSubmitting = false,
}: MovieFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<MovieFormValues>({
    resolver: zodResolver(schema),
    defaultValues: initial
      ? {
          title: initial.title,
          description: initial.description ?? '',
          genre: initial.genre,
          duration: initial.duration,
          rating: initial.rating,
        }
      : {
          title: '',
          description: '',
          genre: MOVIE_GENRES[0],
          duration: 90,
          rating: 7,
        },
  })

  return (
    <form
      onSubmit={handleSubmit(async (values) => {
        await onSubmit(values)
      })}
      className="flex flex-col gap-4 rounded-lg border border-border bg-card p-4"
    >
      <Input
        id="movie-title"
        label="Title"
        {...register('title')}
        error={errors.title?.message}
      />
      <div className="flex flex-col gap-1">
        <label htmlFor="movie-description" className="text-body font-medium text-white">
          Description
        </label>
        <textarea
          id="movie-description"
          rows={3}
          className="rounded-lg border border-border2 bg-card2 px-3 py-2 text-body text-white outline-none transition-colors placeholder:text-muted focus:border-accent"
          {...register('description')}
        />
      </div>
      <Select id="movie-genre" label="Genre" {...register('genre')} error={errors.genre?.message}>
        {MOVIE_GENRES.map((g) => (
          <option key={g} value={g}>
            {g}
          </option>
        ))}
      </Select>
      <Input
        id="movie-duration"
        label="Duration (minutes)"
        type="number"
        min={1}
        {...register('duration', { valueAsNumber: true })}
        error={errors.duration?.message}
      />
      <Input
        id="movie-rating"
        label="Rating (0–10)"
        type="number"
        step="0.1"
        min={0}
        max={10}
        {...register('rating', { valueAsNumber: true })}
        error={errors.rating?.message}
      />
      <div className="flex flex-wrap gap-2">
        <Button type="submit" disabled={isSubmitting} className="inline-flex gap-2">
          {isSubmitting ? <Spinner tone="onPrimary" /> : null}
          {submitLabel}
        </Button>
        {onCancel ? (
          <Button type="button" variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        ) : null}
      </div>
    </form>
  )
}
