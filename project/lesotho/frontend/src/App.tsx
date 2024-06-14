import { Button } from "@/components/ui/button"
import { ModeToggle } from "@/components/mode-toggle"
import { Input } from "./components/ui/input"
import { X } from "lucide-react"
import { useState } from "react"
import { aclCheck, aclUpdate } from "./http/endpoints"
import { useToast } from "./components/ui/use-toast"

function App() {
  const [namespace, setNamespace] = useState("")
  const [object, setObject] = useState("")
  const [relation, setRelation] = useState("")
  const [user, setUser] = useState("")

  const { toast } = useToast()

  const emptyInputs = () => {
    setNamespace("")
    setObject("")
    setRelation("")
    setUser("")
  }

  const checkACL = () => {
    aclCheck({object: `${namespace}:${object}`, relation, user})
      .then(res => toast({
        variant: res.data.authorized ? "default" : "destructive",
        description: `Authorization ${res.data.authorized ? "succeeded" : "failed"}`,
      }))
      // TODO: typed errors shown in toast
      .catch(err => console.log(err))
  }

  const updateACL = () => {
    aclUpdate({object: `${namespace}:${object}`, relation, user})
      .then(() => toast({
        description: "Update succeeded",
      }))
      .catch(err => console.log(err))
  }

  return (
    <>
      <div className="absolute top-[10px] right-[10px]">
        <ModeToggle/>
      </div>
      <div className="m-auto mt-[250px] w-[800px] rounded-lg border">
        <div className="p-4 flex flex-col gap-2">
          <div className="flex gap-2">
            <div className="flex">
              <Input value={namespace} onChange={e => setNamespace(e.target.value)} placeholder="namespace" className="w-[100px] rounded-r-none focus:z-10 focus:rounded-md"></Input>
              <Input value={object} onChange={e => setObject(e.target.value)} placeholder="object" className="w-[150px] rounded-none focus:z-10 focus:rounded-md"></Input>
              <Input value={relation} onChange={e => setRelation(e.target.value)} placeholder="relation" className="w-[150px] rounded-none focus:z-10 focus:rounded-md"></Input>
              <Input value={user} onChange={e => setUser(e.target.value)} placeholder="user" className="w-[150px] rounded-l-none focus:z-10 focus:rounded-md"></Input>
            </div>
            <Button variant="outline" size="icon" onClick={emptyInputs}><X/></Button>
          </div>
          <div className="flex gap-1">
            <Button onClick={() => checkACL()}>Check</Button>
            <Button onClick={() => updateACL()}>Update</Button>
          </div>
        </div>
      </div>
    </>
  )
}

export default App
