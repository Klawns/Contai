import type { ButtonHTMLAttributes } from 'react'

type AuthSubmitButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  isLoading?: boolean
  loadingLabel: string
}

export function AuthSubmitButton({
  children,
  isLoading = false,
  loadingLabel,
  disabled,
  ...buttonProps
}: AuthSubmitButtonProps) {
  return (
    <button
      type="submit"
      className="mt-2 h-11 cursor-pointer rounded-[10px] bg-[#6a22e5] px-4 text-[14px] font-semibold text-white transition-colors hover:bg-[#5a1ec2] disabled:cursor-not-allowed disabled:bg-[#b8a8d8] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
      disabled={disabled || isLoading}
      {...buttonProps}
    >
      {isLoading ? loadingLabel : children}
    </button>
  )
}
