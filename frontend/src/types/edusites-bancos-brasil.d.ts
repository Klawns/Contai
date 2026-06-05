declare module '@edusites/bancos-brasil' {
  export type BancoFormato = 'circulo' | 'quadrado' | 'rounded' | 'sem'

  export type SvgBancoOptions = {
    nome: string
    cor?: string
    fundo?: string
    formato?: BancoFormato
    tamanho?: number
    className?: string
  }

  export function svgBanco(options: SvgBancoOptions): string | null
  export function listarBancos(): string[]
  export function obterPreset(nome: string): unknown
}
