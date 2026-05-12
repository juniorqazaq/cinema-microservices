import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Button } from '../ui/Button'
import { Input } from '../ui/Input'
import { Spinner } from '../ui/Spinner'

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
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
    defaultValues: { name: '', capacity: 100 },
  })

  return (
    <form
      onSubmit={handleSubmit(async (values) => {
        await onSubmit(values)
        reset({ name: '', capacity: 100 })
      })}
      className="flex flex-col gap-3 rounded-lg border border-border bg-card p-4"
    >
      <p className="text-body font-medium text-white">Create hall</p>
      <Input
        id="hall-name"
        label="Name"
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
