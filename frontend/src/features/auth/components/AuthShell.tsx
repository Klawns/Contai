import type { ReactNode } from 'react'

type AuthShellProps = {
  title: string
  subtitle: string
  children: ReactNode
}

export function AuthShell({ title, subtitle, children }: AuthShellProps) {
  return (
    <main className="grid min-h-svh place-items-center bg-[#f4f7fb] px-4 py-8 text-[#241a30]">
      <section className="w-full max-w-[420px] rounded-[18px] border border-[#e5e0ec] bg-white px-5 py-6 shadow-[0_20px_50px_rgba(48,39,61,0.08)] sm:px-7 sm:py-8">
        <div className="mb-6">
          <p className="text-[13px] font-semibold text-[#6a22e5]">Contai</p>
          <h1 className="mt-2 text-[28px] font-semibold leading-tight tracking-0 text-[#241a30]">
            {title}
          </h1>
          <p className="mt-2 text-[14px] leading-6 text-[#6f6679]">{subtitle}</p>
        </div>

        {children}
      </section>
    </main>
  )
}
