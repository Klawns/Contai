import { CalendarDays } from 'lucide-react'
import type { CommitmentType } from '../../types/commitments.ts'

export function EmptyCommitmentState({ type }: { type: CommitmentType }) {
  return (
    <div className="grid flex-1 place-items-center px-5 py-12 text-center">
      <div className="grid justify-items-center">
        <div className="grid h-20 w-20 place-items-center rounded-2xl bg-[#f4f1f7] text-[#6818e8]">
          <CalendarDays className="h-9 w-9" aria-hidden="true" />
        </div>
        <h2 className="mt-4 max-w-[280px] text-[17px] font-semibold leading-snug text-[#2c2237]">
          Nenhuma {type === 'payable' ? 'conta a pagar' : 'conta a receber'} neste mes
        </h2>
        <p className="mt-2 max-w-[300px] text-[13px] font-medium leading-relaxed text-[#81788c]">
          Cadastre compromissos futuros para acompanhar vencimentos sem mexer no saldo.
        </p>
      </div>
    </div>
  )
}
