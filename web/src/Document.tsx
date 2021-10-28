import React from 'react';
export interface DocumentData {
  URL: string
  content: string
  nombre?: string
}

function Document(props: React.PropsWithChildren<DocumentData>) {
  return (<p>{props.nombre || 'Documento sin título'}</p>)
}
export default Document
