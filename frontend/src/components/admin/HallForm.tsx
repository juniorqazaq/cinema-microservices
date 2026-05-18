import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { CITIES, DEFAULT_CITY } from '../../constants/cities'
import { Button } from '../ui/Button'
import { Input, Select } from '../ui/Input'
import { Spinner } from '../ui/Spinner'

const schema = z.object({
  city: z.string().min(1),
  cinema_name: z.string().min(1, 'Cinema name is required'),
  name: z.string().min(1, 'Hall name is required'),
  capacity: z.number().min(1, 'Capacity must be at least 1'),
})

export type HallFormValues = z.infer<typeof schema>

interface HallFormProps {
  onSubmit: (values: HallFormValues) => Promise<void>
  isSubmitting?: boolean
}

export function HallForm({ onSubmit, isSubmitting = false }: HallFormProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<HallFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      city: DEFAULT_CITY,
      cinema_name: '',
      name: '',
      capacity: 100,
    },
  })

  return (
    <form
      onSubmit={handleSubmit(async (values) => {
        await onSubmit(values)
        reset({
          city: DEFAULT_CITY,
          cinema_name: '',
          name: '',
          capacity: 100,
        })
      })}
      className="flex flex-col gap-3 rounded-lg border border-border bg-card p-4"
    >
      <p className="text-body font-medium text-white">Create hall</p>
      <Select id="hall-city" label="City" {...register('city')} error={errors.city?.message}>
        {CITIES.map((c) => (
          <option key={c} value={c}>
            {c}
          </option>
        ))}
      </Select>
      <Input
        id="hall-cinema"
        label="Cinema"
        {...register('cinema_name')}
        error={errors.cinema_name?.message}
      />
      <Input
        id="hall-name"
        label="Hall"
        {...register('name')}
        error={errors.name?.message}
      />
      <Input
        id="hall-capacity"
        label="Capacity"
        type="number"
        min={1}
        {...register('capacity', { valueAsNumber: true })}
        error={errors.capacity?.message}
      />
      <Button type="submit" disabled={isSubmitting} className="w-fit gap-2">
        {isSubmitting ? <Spinner tone="onPrimary" /> : null}
        Create hall
      </Button>
    </form>
  )
}
