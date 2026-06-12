import { ListFilter } from 'lucide-react'

type ReportAllOptionButtonProps = {
  label: string
  isSelected: boolean
  onSelect: () => void
}

export function ReportAllOptionButton({
  label,
  isSelected,
  onSelect,
}: ReportAllOptionButtonProps) {
  return (
    <button
      type="button"
      className={`flex w-full cursor-pointer items-center gap-3 rounded-lg border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
        isSelected ? 'border-[#7b2cff] bg-[#f7f2ff]' : 'border-[#eee8f3] bg-white hover:bg-[#fbf9fe]'
      }`}
      onClick={onSelect}
    >
      <span className="grid h-10 w-10 flex-none place-items-center rounded-full bg-[#f4f1f7] text-[#7b2cff]">
        <ListFilter className="h-5 w-5" aria-hidden="true" />
      </span>
      <span className="min-w-0 flex-1 truncate text-[14px] font-semibold text-[#2c2237]">
        {label}
      </span>
    </button>
  )
}
