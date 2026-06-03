import { LogOut } from 'lucide-react'

type LogoutActionButtonProps = {
  isLoggingOut: boolean
  onLogout: () => void
  variant?: 'sidebar' | 'page'
}

const baseClasses =
  'inline-flex cursor-pointer items-center gap-2 border font-[inherit] font-semibold text-[#c93434] transition-colors hover:border-[#f0caca] hover:bg-[#fff4f4] hover:text-[#a92828] disabled:cursor-not-allowed disabled:text-[#d58d8d] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#d53b3b]'

const variantClasses = {
  sidebar:
    'h-11 w-full justify-start rounded-lg border-transparent bg-transparent px-2.5 text-[13px]',
  page: 'h-11 justify-center rounded-[10px] border-[#f0caca] bg-white px-4 text-[14px] shadow-[0_8px_18px_rgba(48,39,61,0.04)]',
}

export function LogoutActionButton({
  isLoggingOut,
  onLogout,
  variant = 'page',
}: LogoutActionButtonProps) {
  return (
    <button
      type="button"
      className={`${baseClasses} ${variantClasses[variant]}`}
      disabled={isLoggingOut}
      onClick={onLogout}
    >
      <LogOut className="h-4 w-4" aria-hidden="true" />
      {isLoggingOut ? 'Saindo...' : 'Sair'}
    </button>
  )
}
