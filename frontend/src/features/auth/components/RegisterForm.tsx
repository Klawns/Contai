import { zodResolver } from '@hookform/resolvers/zod'
import { AxiosError } from 'axios'
import { useForm } from 'react-hook-form'
import { useRegister } from '../hooks/useRegister.ts'
import { registerSchema } from '../schemas/auth.ts'
import { navigateTo } from '../services/navigation.ts'
import type { RegisterPayload } from '../types/auth.ts'
import { AuthFormField } from './AuthFormField.tsx'
import { AuthSubmitButton } from './AuthSubmitButton.tsx'

function getErrorMessage(error: unknown) {
  if (error instanceof AxiosError && error.response?.status === 409) {
    return 'Ja existe uma conta com este e-mail.'
  }

  return 'Nao foi possivel criar sua conta. Tente novamente.'
}

export function RegisterForm() {
  const registerMutation = useRegister()
  const {
    formState: { errors },
    handleSubmit,
    register,
  } = useForm<RegisterPayload>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
    },
  })

  const onSubmit = handleSubmit((payload) => {
    registerMutation.mutate(payload, {
      onSuccess: () => navigateTo('/', { replace: true }),
    })
  })

  return (
    <form className="grid gap-4" onSubmit={onSubmit}>
      <AuthFormField
        label="Nome"
        type="text"
        autoComplete="name"
        error={errors.name}
        {...register('name')}
      />
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
        autoComplete="new-password"
        error={errors.password}
        {...register('password')}
      />

      {registerMutation.isError ? (
        <p className="rounded-[10px] border border-[#f0caca] bg-[#fff8f8] px-3 py-2 text-[13px] font-medium text-[#b93838]">
          {getErrorMessage(registerMutation.error)}
        </p>
      ) : null}

      <AuthSubmitButton isLoading={registerMutation.isPending} loadingLabel="Criando conta...">
        Criar conta
      </AuthSubmitButton>

      <button
        type="button"
        className="cursor-pointer text-[13px] font-semibold text-[#6a22e5] hover:text-[#5a1ec2] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
        onClick={() => navigateTo('/login')}
      >
        Ja tenho conta
      </button>
    </form>
  )
}
