import axios from 'axios'

export function mapApiError(status: number | undefined): string | null {
  if (status === 400) return 'Please check your input and try again.'
  if (status === 403) return "You don't have permission to do this."
  if (status === 404) return 'Not found'
  if (status === 409) return 'This seat was just taken. Please choose another.'
  if (status === 500) return 'Something went wrong. Please try again.'
  return null
}

export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const body = error.response?.data as { error?: string } | undefined
    if (body?.error && typeof body.error === 'string') {
      return body.error
    }
    const mapped = mapApiError(error.response?.status)
    if (mapped) return mapped
  }
  if (error instanceof Error) {
    return error.message
  }
  return 'Something went wrong. Please try again.'
}
