import { Link, useParams } from 'react-router-dom'

const buyTicketBtnClass =
  'inline-flex min-w-[140px] items-center justify-center rounded-lg border border-accent bg-accent px-4 py-2.5 text-body font-semibold text-white shadow-sm transition-colors hover:border-accentHover hover:bg-accentHover'

export function BuyTicketButton() {
  const { id } = useParams()
  const target = id ? `/movies/${id}/sessions` : '/movies'

  return (
    <Link to={target} className={buyTicketBtnClass}>
      Buy ticket
    </Link>
  )
}
