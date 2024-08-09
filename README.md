## Chạy chương trình:

Ta cũng sẽ cần tạo một tài khoản IAM trên AWS và dùng access_key, secret_access_key để sử dụng chức năng S3

Ta sẽ cần tạo 1 file .env chứa các dòng sau:

![image](https://github.com/user-attachments/assets/d020ebe2-e099-40fc-8f54-0b5144e62e0e)

Chạy project:
```bash
cd cmd 
air
```

## Mô tả:
App Todo BeBetter là một ứng dụng ghi chú giúp người dùng dễ dàng ghi nhớ và kiểm soát các công việc cần và đã thực hiện của mình. Với giao diện bắt mắt, tối giản, người dùng có thể dễ dàng làm quen và sử dụng ứng dụng một cách nhanh chóng

### App bao gồm các chức năng chính như:
- Tạo tài khoản, đăng nhập
- Cập nhật thông tin tài khoản, user
- Thêm, xóa, sửa các ghi chú
- Tìm kiếm ghi chú theo ngày
- Cập nhật trạng thái ghi chứ

### Công nghệ đã sử dụng
- [Golang(Echo framework)](https://echo.labstack.com/): Thiết lập các api
- [JWT](https://jwt.io/): Bảo mật tài khoản user
- [Sqlx](https://github.com/jmoiron/sqlx): Library hỗ trợ kết nối database và thực hiện viết các câu query để lấy dữ liệu database
- AWS S3: Lưu trữ dữ liệu avatar user
- Databse: MySQL

## FE Proeject
https://github.com/YonkLongSchlong/RN-Todo

## Link video demo:
https://drive.google.com/file/d/1fbFF7c38f2EtaK-jej8gGdcZEZZtHYuJ/view?usp=drive_link
  



