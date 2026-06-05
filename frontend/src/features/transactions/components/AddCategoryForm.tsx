import { Controller, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useCreateCategory } from '../hooks/useCreateCategory.ts'
import { createCategoryPayloadSchema } from '../schemas/transactions.ts'
import type { CategoryTransactionType } from '../types/transactions.ts'
import { FormActionButton, FormSection } from './FormFields.tsx'
import { SelectionSheet } from './SelectionSheet.tsx'

type AddCategoryFormValues = {
  name: string
  color: string
  icon: string
}

type AddCategoryFormProps = {
  type: CategoryTransactionType
  isOpen: boolean
  onClose: () => void
  onCreated: (categoryId: string) => void
}

const defaultValues: AddCategoryFormValues = {
  name: '',
  color: '#7B2CFF',
  icon: 'tag',
}

export function AddCategoryForm({
  type,
  isOpen,
  onClose,
  onCreated,
}: AddCategoryFormProps) {
  const createCategoryMutation = useCreateCategory(type)
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AddCategoryFormValues>({
    defaultValues,
    resolver: zodResolver(createCategoryPayloadSchema.omit({ type: true })),
  })

  return (
    <SelectionSheet title="Nova categoria" isOpen={isOpen} onClose={onClose}>
      <form
        className="grid gap-4"
        onSubmit={handleSubmit((values) => {
          createCategoryMutation.mutate(values, {
            onSuccess: (category) => {
              reset(defaultValues)
              onCreated(category.id)
              onClose()
            },
          })
        })}
      >
        <Controller
          control={control}
          name="name"
          render={({ field }) => (
            <FormSection label="Nome" error={errors.name?.message}>
              <input
                className="h-12 w-full rounded-lg border border-[#e5deee] bg-white px-4 text-[15px] font-medium text-[#2c2237] outline-none focus:border-[#7b2cff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                value={field.value}
                onChange={field.onChange}
              />
            </FormSection>
          )}
        />
        <Controller
          control={control}
          name="color"
          render={({ field }) => (
            <FormSection label="Cor" error={errors.color?.message}>
              <input
                type="color"
                className="h-12 w-full cursor-pointer rounded-lg border border-[#e5deee] bg-white px-3 py-2"
                value={field.value}
                onChange={field.onChange}
              />
            </FormSection>
          )}
        />
        <Controller
          control={control}
          name="icon"
          render={({ field }) => (
            <FormSection label="Icone" error={errors.icon?.message}>
              <input
                className="h-12 w-full rounded-lg border border-[#e5deee] bg-white px-4 text-[15px] font-medium text-[#2c2237] outline-none focus:border-[#7b2cff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7b2cff]"
                value={field.value}
                onChange={field.onChange}
              />
            </FormSection>
          )}
        />
        {createCategoryMutation.isError ? (
          <p className="rounded-lg border border-[#f0caca] bg-[#fff7f7] px-3 py-2 text-[13px] font-medium text-[#b93838]">
            Nao foi possivel criar a categoria.
          </p>
        ) : null}
        <FormActionButton isPending={createCategoryMutation.isPending}>
          Salvar categoria
        </FormActionButton>
      </form>
    </SelectionSheet>
  )
}
