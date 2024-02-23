import {useSelector} from 'react-redux';

export default function Status() {
  const data=useSelector(state=>state.data);

  return (
    <div>
      {JSON.stringify(data)}
    </div>
  );
}