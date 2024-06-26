import { Button } from "@/components/ui/button"
import { ModeToggle } from "@/components/mode-toggle"
import { Input } from "./components/ui/input"
import { X } from "lucide-react"
import { useRef, useState } from "react"
import { aclCheck, aclUpdate, namespaceUpdate } from "./http/endpoints"
import { useToast } from "./components/ui/use-toast"
import axios, { AxiosError } from "axios"
import { Separator } from "./components/ui/separator"

function App() {
  const [namespace, setNamespace] = useState("")
  const [object, setObject] = useState("")
  const [relation, setRelation] = useState("")
  const [user, setUser] = useState("")
  const [aclLocked, setAclLocked] = useState(false);
  const aclAnyFieldEmpty = namespace == "" || object == "" || relation == "" || user == ""
  const aclAllFieldEmpty = namespace == "" && object == "" && relation == "" && user == ""

  const [namespaceFile, setNamespaceFile] = useState<File | null>(null)
  const [namespaceLocked, setNamespaceLocked] = useState(false)
  const namespaceInputElement = useRef<HTMLInputElement | null>(null)
  const namespaceEmpty = namespaceFile == null

  const { toast } = useToast()

  const emptyInputs = () => {
    setNamespace("")
    setObject("")
    setRelation("")
    setUser("")
  }

  const checkACL = () => {
    setAclLocked(true)
    const namespaceObject = `${namespace}:${object}`
    aclCheck({object: namespaceObject, relation, user})
      .then(res => toast({
        title: "ACL Check result",
        description: `${res.data.authorized ? "Yes" : "No"}, ${user} -> ${relation} of ${namespaceObject}`,
      }))
      // https://github.com/axios/axios/issues/3612#issuecomment-770224236
      .catch((err: Error | AxiosError) => {
        if (axios.isAxiosError(err) && err.response?.status === 400) {
          toast({
            title: "ACL Check failed",
            variant: "destructive",
            description: err.response.data as string
          })
        } else {
          toast({
            title: "Unexpected error",
            variant: "destructive",
            description: "Please check your console"
          })
          console.log(err.message)
        }
      })
      .finally(() => setAclLocked(false))
  }

  const updateACL = () => {
    setAclLocked(true)
    const namespaceObject = `${namespace}:${object}`
    aclUpdate({object: namespaceObject, relation, user})
      .then(() => toast({
        title: "ACL Update result",
        description: `Successfully updated ${namespaceObject}#${relation}@${user}`,
      }))
      .catch((err: Error | AxiosError) => {
        if (axios.isAxiosError(err) && err.response?.status === 400) {
          toast({
            title: "ACL Update failed",
            variant: "destructive",
            description: err.response.data as string
          })
        } else {
          toast({
            title: "Unexpected error",
            variant: "destructive",
            description: "Please check your console"
          })
          console.log(err.message)
        }
      })
      .finally(() => setAclLocked(false))
  }

  const attemptStringtoJSON = (val: string): unknown => {
    try {
      return JSON.parse(val)
    } catch {
      return null
    }
  }

  const updateNamespace = () => {
    if (namespaceFile === null) {
      return
    }
    setNamespaceLocked(true)
    const reader = new FileReader()
    reader.onload = () => {
      const namespace = reader.result as string
      const json = attemptStringtoJSON(namespace)
      if (json === null) {
        toast({
          title: "Namespace Update failed",
          variant: "destructive",
          description: `${namespaceFile.name} is not a valid JSON file`,
        })
        setNamespaceLocked(false)
        return
      }
      namespaceUpdate(json)
        .then(() => {
          toast({
            title: "Namespace Update result",
            description: "Successfully updated",
          })
          setNamespaceFile(null)
          if (namespaceInputElement.current) {
            namespaceInputElement.current.value = ""
          }
        })
        .catch((err: Error | AxiosError) => {
          if (axios.isAxiosError(err) && err.response?.status === 400) {
            toast({
              title: "Namespace Update failed",
              variant: "destructive",
              description: err.response.data as string
            })
          } else {
            toast({
              title: "Unexpected error",
              variant: "destructive",
              description: "Please check your console"
            })
            console.log(err.message)
          }
        })
        .finally(() => {
          setNamespaceLocked(false)
        })
    }
    reader.readAsText(namespaceFile)
  }

  return (
    <>
      <div className="flex justify-end p-4">
        <ModeToggle/>
      </div>
      <div className="m-auto mt-[250px] w-[800px] rounded-lg border">
        <div className="p-4 flex flex-col gap-2">
          <h2>ACL</h2>
          <div className="flex gap-2">
            <div className="flex">
              <Input disabled={aclLocked} value={namespace} onChange={e => setNamespace(e.target.value)} placeholder="namespace" className="w-[100px] rounded-r-none focus:z-10 focus:rounded-md"></Input>
              <Input disabled={aclLocked} value={object} onChange={e => setObject(e.target.value)} placeholder="object" className="w-[150px] rounded-none focus:z-10 focus:rounded-md"></Input>
              <Input disabled={aclLocked} value={relation} onChange={e => setRelation(e.target.value)} placeholder="relation" className="w-[150px] rounded-none focus:z-10 focus:rounded-md"></Input>
              <Input disabled={aclLocked} value={user} onChange={e => setUser(e.target.value)} placeholder="user" className="w-[150px] rounded-l-none focus:z-10 focus:rounded-md"></Input>
            </div>
            <Button disabled={aclLocked || aclAllFieldEmpty} variant="outline" size="icon" onClick={emptyInputs}><X/></Button>
          </div>
          <div className="flex gap-1">
            <Button disabled={aclLocked || aclAnyFieldEmpty} onClick={() => checkACL()}>Check</Button>
            <Button disabled={aclLocked || aclAnyFieldEmpty} onClick={() => updateACL()}>Update</Button>
          </div>
        </div>
        <Separator/>
        <div className="p-4 flex flex-col gap-2">
          <h2>Namespace</h2>
          <div className="flex gap-2">
            <Input ref={namespaceInputElement} disabled={namespaceLocked} type="file" accept=".json" className="w-[300px]" onChange={e => setNamespaceFile(e.target.files?.item(0) ?? null)}/>
            <Button disabled={namespaceLocked || namespaceEmpty} onClick={() => updateNamespace()}>Update</Button>
          </div>
        </div>
      </div>
    </>
  )
}

export default App
