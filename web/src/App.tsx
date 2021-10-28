import React, { useState } from 'react';
import searchImg from './images/system-search-symbolic.svg';
import './App.css';
import Box from './Box';
import Section from './Section';
import Button from './Button';

import Expediente, { ExpedienteData } from './Expediente';

interface ExpedienteURL {
  id: string
  file: string
  selected: boolean
}

interface ExpDataMap {
  [id: string]: ExpedienteData
}

function App() {
  const [search, setSearch] = useState("");
  const [expedientes, setExpedientes] = useState<ExpedienteURL[]>([]);
  const [data, setData] = useState<ExpDataMap>();
  const [selected, setSelected] = useState<ExpedienteData>();

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    event.preventDefault();
    setSearch(event.target.value)
  }

  const handleSelect = (e: ExpedienteData) => (event: React.MouseEvent<HTMLAnchorElement>) => {
    event.preventDefault();
    setSelected(e)
  }


  useEffect(() => {
    (async () => {
      const res = await fetch(`${process.env.PUBLIC_URL}/data/expedientes.json`);
      const json = await res.json();
      setExpedientes(json)
    })();
  }, [])

  useEffect(() => {
    if (!expedientes.length) return;
    (async () => {
      for (let i = 0; i < expedientes.length; i++) {
        const res = await fetch(`${process.env.PUBLIC_URL}/data/${expedientes[i].file}.json`)
        const json = await res.json();
        setData({ ...data, [json.ExpId]: json })
      }
    })()
  }, [expedientes, data])

  return (
    <div className="App">
      <Box id="title">
        <h1>Buscador de Actuaciones de O.D.I.A.</h1>
      </Box>
      <Section title="expedientes">
        {(expedientes.length && data)
          ? expedientes.map(({ id }, i) =>
            <Button key={i}
              onClick={handleSelect(data[id])}
              selected={selected === data[i]}>
              {id}
            </Button>)
          : <p>loading...</p>}
        {selected !== undefined && <Expediente {...selected} />}
      </Section>
      <div title="buscar">
        <input type="search" id="search-input" name="search"
          onChange={handleChange} value={search} placeholder="buscar..." />
        <img src={searchImg} alt="Buscar" />
        {search}
      </div>
    </div >
  );
}

export default App;
