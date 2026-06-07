import { typeCopy } from '../../lib/commitmentPresentation.ts'
import type { CommitmentType } from '../../types/commitments.ts'

export function PlanningTabs({
  type,
  onSelectType,
}: {
  type: CommitmentType
  onSelectType: (type: CommitmentType) => void
}) {
  return (
    <div className="grid grid-cols-2 rounded-full bg-[#f4f1f7] p-1">
      {(['payable', 'receivable'] as const).map((option) => (
        <button
          key={option}
          type="button"
          className={`h-10 cursor-pointer rounded-full text-[13px] font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff] ${
            type === option
              ? 'bg-white text-[#2c2237] shadow-[0_4px_12px_rgba(40,29,53,0.08)]'
              : 'text-[#6f647b] hover:text-[#2c2237]'
          }`}
          onClick={() => onSelectType(option)}
        >
          {typeCopy[option].plural}
        </button>
      ))}
    </div>
  )
}
