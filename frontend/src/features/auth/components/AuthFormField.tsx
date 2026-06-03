import type { ComponentPropsWithoutRef } from 'react'
import type { FieldError } from 'react-hook-form'

type AuthFormFieldProps = ComponentPropsWithoutRef<'input'> & {
  label: string
  error?: FieldError
}

export function AuthFormField({ label, error, id, ...inputProps }: AuthFormFieldProps) {
  const fieldId = id ?? inputProps.name
  const errorId = error ? `${fieldId}-error` : undefined

  return (
    <label className="grid gap-2 text-[13px] font-medium text-[#4b4355]" htmlFor={fieldId}>
      {label}
      <input
        id={fieldId}
        aria-invalid={Boolean(error)}
        aria-describedby={errorId}
        className="h-11 rounded-[10px] border border-[#dcd6e6] bg-white px-3 text-[14px] text-[#241a30] outline-none transition-colors placeholder:text-[#9f97aa] focus:border-[#7b2cff] focus:ring-2 focus:ring-[#7b2cff]/15"
        {...inputProps}
      />
      {error ? (
        <span id={errorId} className="text-[12px] font-medium text-[#b93838]">
          {error.message}
        </span>
      ) : null}
    </label>
  )
}
