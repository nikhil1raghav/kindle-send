from email.message import EmailMessage
from config import sender, password
import mimetypes
import smtplib
import os
def sendMail(receiver, file_path):
    """
    Sends file to receiver
    from email of sender
    credentials defined in config
    """
    msg = EmailMessage()

    msg['To']=receiver
    msg['From']=sender

    file_name = os.path.basename(file_path)

    with open(file_path, "rb") as f:
        file_data = f.read()
    mime_type, _ = mimetypes.guess_type(file_name)
    if mime_type:
        mime_type, mime_subtype = mime_type.split("/",1)
    else:
        mime_type = "application"
        mime_subtype = "octet-stream"
    
    msg['Subject'] = file_name
    msg.add_attachment(file_data, maintype = mime_type, subtype = mime_subtype, filename = file_name)
    
    with smtplib.SMTP_SSL('smtp.gmail.com', 465) as smtp:
        smtp.login(sender, password)
        smtp.send_message(msg)

    print("Sent ", file_name, "to", receiver)
