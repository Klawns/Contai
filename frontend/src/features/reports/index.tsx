import { DashboardLayout } from '../dashboard/components'
import { ReportsFilterForm } from './components/ReportsFilterForm.tsx'
import { ReportsPageLayout } from './components/ReportsPageLayout.tsx'
import { useReportFiltersForm } from './hooks/useReportFiltersForm.ts'

export function ReportsPage() {
  const {
    formState,
    updateField,
    isExporting,
    message,
    handleSubmit,
  } = useReportFiltersForm()

  return (
    <DashboardLayout width="full">
      <ReportsPageLayout>
        <ReportsFilterForm
          formState={formState}
          isExporting={isExporting}
          message={message}
          onFieldChange={updateField}
          onSubmit={handleSubmit}
        />
      </ReportsPageLayout>
    </DashboardLayout>
  )
}
