import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import './App.css'
import App from './App.tsx'

// Add global error logging for debugging
window.addEventListener('error', (event) => {
  console.error('🚨 Global error:', event.error)
  console.error('🚨 Error details:', event)
})

window.addEventListener('unhandledrejection', (event) => {
  console.error('🚨 Unhandled promise rejection:', event.reason)
})

console.log('🚀 App starting up...')

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
