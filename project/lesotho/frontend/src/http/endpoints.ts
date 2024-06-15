import axios from "axios";

export interface ACLDirective {
    object: string,
    relation: string,
    user: string,
}

export interface AuthorizationResponse {
    authorized: boolean,
}

const api = axios.create({
    baseURL: "http://localhost:5001/"
})

export const aclCheck = (directive: ACLDirective) => {
    return api.get<AuthorizationResponse>("acl/check", { params: { ...directive } })
}

export const aclUpdate = (directive: ACLDirective) => {
    return api.post<void>("acl", directive)
}

export const namespaceUpdate = (namespace: string) => {
    return api.post<void>("namespace", JSON.parse(namespace))
}
