import { Button } from "@/components/ui/button"
import { ThemeProvider } from "@/components/theme-provider"
import { ModeToggle } from "@/components/mode-toggle"

function App() {
  return (
    <ThemeProvider>
      <ModeToggle/>
      <Button>Click me</Button>
    </ThemeProvider>
  )
}

export default App
