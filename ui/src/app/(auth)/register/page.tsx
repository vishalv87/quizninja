import { Metadata } from 'next'
import { RegisterForm } from '@/components/auth/RegisterForm'

export const metadata: Metadata = {
  title: 'Register | QuizNinja',
  description: 'Create a new QuizNinja account',
}

export default function RegisterPage() {
  return <RegisterForm />
}