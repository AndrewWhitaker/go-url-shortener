import './App.css';
import {Button, Form} from "react-bootstrap"
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <section>
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>Long URL</Form.Label>
        <Form.Control type="url" placeholder="Enter URL" />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Slug</Form.Label>
        <Form.Control type="url" placeholder="Enter Slug" />
      </Form.Group>
      <Button variant="primary" type="submit">Submit</Button>
    </Form>
    </section>
  );
}

export default App;
