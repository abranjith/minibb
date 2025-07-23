import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'

type Theme = 'light' | 'dark'

interface ThemeContextType {
  theme: Theme
  toggleTheme: () => void
  setTheme: (theme: Theme) => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export const useTheme = () => {
  const context = useContext(ThemeContext)
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider')
  }
  return context
}

interface ThemeProviderProps {
  children: ReactNode
}

export const ThemeProvider = ({ children }: ThemeProviderProps) => {
  // Initialize theme from localStorage or default to 'light'
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window !== 'undefined') {
      const savedTheme = localStorage.getItem('minibb-theme') as Theme
      return savedTheme || 'light'
    }
    return 'light'
  })

  // Apply theme to document when theme changes
  useEffect(() => {
    const root = document.documentElement
    console.log('🎨 Applying theme:', theme)
    if (theme === 'dark') {
      root.setAttribute('data-theme', 'dark')
      console.log('🌙 Dark theme applied, data-theme attribute:', root.getAttribute('data-theme'))
    } else {
      root.removeAttribute('data-theme')
      console.log('☀️ Light theme applied, data-theme attribute removed')
    }
    // Save to localStorage
    localStorage.setItem('minibb-theme', theme)
  }, [theme])

  const toggleTheme = () => {
    const newTheme = theme === 'light' ? 'dark' : 'light'
    console.log('🔄 Toggling theme from', theme, 'to', newTheme)
    setThemeState(newTheme)
  }

  const setTheme = (newTheme: Theme) => {
    setThemeState(newTheme)
  }

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}
