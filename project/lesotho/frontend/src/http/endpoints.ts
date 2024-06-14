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
    baseURL: "http://localhost:5000/"
})

export const aclCheck = (directive: ACLDirective) => {
    return api.get<AuthorizationResponse>("acl/check", { params: {
        object: directive.object,
        relation: directive.relation,
        user: directive.user,
    }})
}

export const aclUpdate = (directive: ACLDirective) => {
    return api.post<void>("acl", directive)
}
