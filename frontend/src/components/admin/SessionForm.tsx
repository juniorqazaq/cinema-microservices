import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { Movie } from '../../types'
import { CITIES, DEFAULT_CITY } from '../../constants/cities'
import { Button } from '../ui/Button'
import { Input, Select } from '../ui/Input'
import { Spinner } from '../ui/Spinner'

const schema = z.object({
  city: z.string().min(1),
  movie_id: z.string().min(1, 'Movie is required'),
  hall_id: z.string().min(1, 'Hall is required'),
  start_time: z.string().min(1, 'Start time is required'),
  price: z.number().min(0, 'Price must be 0 or more'),
})

export type SessionFormValues = z.infer<typeof schema>

interface SessionFormProps {
  movies: Movie[]
  hallOptions: { id: string; name: string; city?: string }[]
  onSubmit: (values: SessionFormValues) => Promise<void>
  isSubmitting?: boolean
}

export function SessionForm({
  movies,
  hallOptions,
  onSubmit,
  isSubmitting = false,
}: SessionFormProps) {
  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<SessionFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      city: DEFAULT_CITY,
      movie_id: movies[0]?.id ?? '',
      hall_id: hallOptions[0]?.id ?? '',
      start_time: '',
      price: 2500,
    },
  })

  const city = watch('city')
  const hallsInCity = hallOptions.filter((h) => !h.city || h.city === city)

  return (
    <form
      onSubmit={handleSubmit(async (values) => {
        await onSubmit(values)
        reset((f) => ({ ...f, start_time: '' }))
      })}
      className="flex flex-col gap-3 rounded-lg border border-border bg-card p-4"
    >
      <p className="text-body font-medium text-white">Create session</p>
      <Select id="session-city" label="City" {...register('city')} error={errors.city?.message}>
        {CITIES.map((c) => (
          <option key={c} value={c}>
            {c}
          </option>
        ))}
      </Select>
      <Select
        id="session-movie"
        label="Movie"
        {...register('movie_id')}
        error={errors.movie_id?.message}
      >
        {movies.map((m) => (
          <option key={m.id} value={m.id}>
            {m.title}
          </option>
        ))}
      </Select>
      <Select
        id="session-hall"
        label="Hall / cinema"
        {...register('hall_id')}
        error={errors.hall_id?.message}
      >
        {hallsInCity.length === 0 ? (
          <option value="">Create a hall in {city} first</option>
        ) : (
          hallsInCity.map((h) => (
            <option key={h.id} value={h.id}>
              {h.name}
            </option>
          ))
        )}
      </Select>
      <Input
        id="session-start"
        label="Start time"
        type="datetime-local"
        {...register('start_time')}
        error={errors.start_time?.message}
      />
      <Input
        id="session-price"
        label="Price (₸)"
        type="number"
        min={0}
        {...register('price', { valueAsNumber: true })}
        error={errors.price?.message}
      />
      <Button type="submit" disabled={isSubmitting} className="w-fit gap-2">
        {isSubmitting ? <Spinner tone="onPrimary" /> : null}
        Create session
      </Button>
    </form>
  )
}
