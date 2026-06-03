import { zodResolver } from '@hookform/resolvers/zod'
import { AxiosError } from 'axios'
import { useForm } from 'react-hook-form'
import { loginSchema } from '../schemas/auth.ts'
import type { LoginPayload } from '../types/auth.ts'
import { useLogin } from '../hooks/useLogin.ts'
import { navigateTo } from '../services/navigation.ts'
import { AuthFormField } from './AuthFormField.tsx'
import { AuthSubmitButton } from './AuthSubmitButton.tsx'

function getErrorMessage(error: unknown) {
  if (error instanceof AxiosError && error.response?.status === 401) {
    return 'E-mail ou senha invalidos.'
  }

  return 'Nao foi possivel entrar. Tente novamente.'
}

export function LoginForm() {
  const loginMutation = useLogin()
  const {
    formState: { errors },
    handleSubmit,
    register,
  } = useForm<LoginPayload>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const onSubmit = handleSubmit((payload) => {
    loginMutation.mutate(payload, {
      onSuccess: () => navigateTo('/', { replace: true }),
    })
  })

  return (
    <form className="grid gap-4" onSubmit={onSubmit}>
      <AuthFormField
        label="E-mail"
        type="email"
        autoComplete="email"
        error={errors.email}
        {...register('email')}
      />
      <AuthFormField
        label="Senha"
        type="password"
        autoComplete="current-password"
        error={errors.password}
        {...register('password')}
      />

      {loginMutation.isError ? (
        <p className="rounded-[10px] border border-[#f0caca] bg-[#fff8f8] px-3 py-2 text-[13px] font-medium text-[#b93838]">
          {getErrorMessage(loginMutation.error)}
        </p>
      ) : null}

      <AuthSubmitButton isLoading={loginMutation.isPending} loadingLabel="Entrando...">
        Entrar
      </AuthSubmitButton>

      <button
        type="button"
        className="cursor-pointer text-[13px] font-semibold text-[#6a22e5] hover:text-[#5a1ec2] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        onClick={() => navigateTo('/registro')}
      >
        Criar conta
      </button>
    </form>
  )
}
