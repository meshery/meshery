import React, { useState, useEffect } from 'react';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import { Paper, TextField } from '@mui/material';
import { SubmitButton, DescriptionModal } from './styles';
import { Formik } from 'formik';
import * as yup from 'yup';
import axios from 'axios';
import { Field, Form } from 'formik';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';
import theme from '../../themes/app';
// import ResponseIcon from '../../assets/icons/support/response';

const SupportForm = (props) => {
  const [submit, setSubmit] = useState(false);
  const { user } = props;
  console.log('user', user);
  const [memberFormOne, setmemberFormOne] = useState({
    firstname: '',
    lastname: '',
    email: '',
    subject: '',
    message: '',
    scope: '',
    form: 'contact',
  });

  const FormSchema = yup.object().shape({
    subject: yup.string().min(5, 'Too short').required('Required'),
    message: yup
      .string()
      .min(5, 'Please offer a longer description of you inquiry')
      .required('Required'),
  });

  const errorMsg = {
    color: 'red',
    textAlign: 'left',
  };

  useEffect(() => {
    if (user) {
      const { first_name, last_name, email } = user;
      setmemberFormOne((prevForm) => ({
        ...prevForm,
        firstname: first_name,
        lastname: last_name,
        email: email,
      }));
    }
    if (submit) {
      axios.post('https://hook.us1.make.com/r5qgpjel5tlhtyndcgjvkrdkoc65417y', {
        memberFormOne: {
          ...memberFormOne,
          firstname: memberFormOne.firstname,
          lastname: memberFormOne.lastname,
          email: memberFormOne.email,
          subject: memberFormOne.subject,
          message: memberFormOne.message,
        },
      });
    }
  }, [
    user,
    submit,
    memberFormOne.email,
    memberFormOne.firstname,
    memberFormOne.lastname,
    memberFormOne.subject,
    memberFormOne.message,
  ]);

  if (submit) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <Paper sx={{ mt: 4, paddingBottom: '3.75rem', maxWidth: 700 }}>
          <Box
            component="span"
            sx={{
              p: 2,
              border: '1px dashed grey',
              width: '70%',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              mt: 8,
              mx: 'auto',
            }}
          >
            <Typography sx={{ fontSize: 48, fontWeight: 500 }}>
              {/* <ResponseIcon width={128} /> */}
            </Typography>
            <Typography sx={{ fontSize: 28, fontWeight: 400, textAlign: 'center', mt: 4 }}>
              Your response has been recorded we will endeavor to promptly contact you with a
              suitable solution.
            </Typography>
          </Box>
        </Paper>
      </div>
    );
  }

  return (
    <>
      {/* <SupportModal sx={{ mt: 2 }}> */}
      <Formik
        initialValues={memberFormOne}
        onSubmit={(values) => {
          setmemberFormOne(values);
          setSubmit(true);
        }}
        validationSchema={FormSchema}
      >
        {({ errors, touched }) => (
          <Form classname="form" method="post">
            <DescriptionModal>
              <Typography sx={{ fontSize: 32, fontWeight: 500, mt: 10 }}>
                Need more help or have a different question?
              </Typography>
              <Typography sx={{ fontSize: 22, fontWeight: 500, mt: 1 }}>
                Send us your inquiry and we’ll get back to you promptly.
              </Typography>
            </DescriptionModal>

            <Box sx={{ width: 700, maxWidth: '100%', mx: 'auto', mt: 6 }}>
              <Field
                as={TextField}
                fullWidth
                htmlFor="subject"
                label="Subject"
                id="subject"
                type="text"
                className="text-field"
                name="subject"
                maxLength="32"
                pattern="[A-Za-z]{1,32}"
                required
                onInvalid={(e) => e.target.setCustomValidity('Please fill-in this field')}
                sx={{
                  '& .MuiInputBase-input': {
                    height: '1.25rem',
                  },
                  mt: 2,
                }}
                onInput={(e) => e.target.setCustomValidity('')}
              />
              {errors.subject && touched.subject ? (
                <div style={errorMsg}>{errors.subject}</div>
              ) : null}
            </Box>
            <Box sx={{ display: 'flex', flexDirection: 'column' }}>
              <Box sx={{ flexGrow: 1, width: 700, maxWidth: '100%', mx: 'auto', mt: 2 }}>
                <Field
                  as={TextField}
                  fullWidth
                  htmlFor="message"
                  label="Issue Description"
                  id="message"
                  type="text"
                  className="text-field"
                  name="message"
                  maxLength="32"
                  pattern="[A-Za-z]{1,32}"
                  required
                  onInvalid={(e) => e.target.setCustomValidity('Please fill-in this field')}
                  sx={{
                    '& .MuiInputBase-input': {
                      height: '5rem',
                    },
                  }}
                  onInput={(e) => e.target.setCustomValidity('')}
                />
                {errors.message && touched.message ? (
                  <div style={errorMsg}>{errors.message}</div>
                ) : null}
                <FormControl sx={{ float: 'left', mt: 4 }}>
                  <FormLabel
                    id="demo-radio-buttons-group-label"
                    sx={{ fontSize: 20, fontWeight: 500 }}
                  >
                    Scope Of Question
                  </FormLabel>
                  <Field name="scope" htmlFor="scope">
                    {({ field }) => (
                      <RadioGroup
                        aria-labelledby="demo-radio-buttons-group-label"
                        {...field}
                        sx={{ mt: 2 }}
                      >
                        <FormControlLabel
                          value="account/billing"
                          name="scope"
                          control={<Radio />}
                          label="Account/Billing"
                        />
                        <FormControlLabel
                          value="sales"
                          name="scope"
                          control={<Radio />}
                          label="Sales"
                        />
                        <FormControlLabel
                          value="support"
                          name="scope"
                          control={<Radio />}
                          label="Support"
                        />
                        <FormControlLabel
                          value="partnership"
                          name="scope"
                          control={<Radio />}
                          label="Partnership"
                        />
                        {/* <FormControlLabel
                            value="landscape"
                            name="scope"
                            control={<Radio />}
                            label="Extension Points"
                          /> */}
                        <FormControlLabel
                          value="community"
                          name="scope"
                          control={<Radio />}
                          label="Community"
                        />
                      </RadioGroup>
                    )}
                  </Field>
                </FormControl>
              </Box>
              <SubmitButton
                variant="contained"
                sx={{ mt: 2, width: 8, mx: 'auto', background: theme.palette.secondary.focused }}
                type="submit"
              >
                Send
              </SubmitButton>
            </Box>
          </Form>
        )}
      </Formik>
      {/* </SupportModal> */}
    </>
  );
};

export default SupportForm;
